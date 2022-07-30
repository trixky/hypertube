import { stat, writeFile } from 'fs/promises';
import { resolve } from 'path';
import { Router } from 'express';
import jdenticon from 'jdenticon';
import { CACHE_PATH } from '../lib/storage';
import { getUser, setUserExtension } from '../postgres/users';

const router = Router();

function randomPicture(id: number | string) {
	return jdenticon.toPng(id, 128);
}

router.get('/users/:userId/picture', async function (req, res) {
	// Sanitize userId
	const userId = parseInt(req.params.userId);
	if (isNaN(userId) || userId < 1) {
		return res.type('image/png').status(400).send(randomPicture(userId));
	}

	// ... check that the user exists
	const userQueryResult = await getUser(userId);
	if (!userQueryResult || userQueryResult.rows.length < 1) {
		return res.type('image/png').status(404).send(randomPicture(userId));
	}

	// If the user already has a picture
	// -- check if it exists and return it
	const user = userQueryResult.rows[0];
	if (user.extension) {
		try {
			const path = `${CACHE_PATH}/${userId}.${user.extension}`;
			const absolutePath = resolve(path);
			await stat(absolutePath);
			return res.status(200).type(`image/${user.extension}`).sendFile(absolutePath);
		} catch (error) {
			console.error(error);
			await setUserExtension(userId, null);
		}
	}

	// -- else check if the random picture is already generated
	const path = `${CACHE_PATH}/${userId}.png`;
	const absolutePath = resolve(path);
	try {
		await stat(absolutePath);
		return res.type('image/png').status(200).sendFile(absolutePath);
	} catch (err) {
		// A new picture will be generated
	}

	// -- else generate a picture and save it to disk
	const picture = randomPicture(userId);
	try {
		await writeFile(absolutePath, picture);
	} catch (error) {
		// Failed to save to disk, but the image can still be returned
		console.error(error);
	}
	return res.type('image/png').status(200).send(picture);
});

export default router;
