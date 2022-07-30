import { stat, unlink, writeFile } from 'fs/promises';
import { resolve } from 'path';
import { Router } from 'express';
import jdenticon from 'jdenticon';
import { CACHE_PATH } from '../lib/storage';
import { getUser, setUserExtension } from '../postgres/users';

const router = Router();

function randomPicture(id: number | string) {
	return jdenticon.toPng(id, 128);
}

const allowedExtensions = ['png', 'jpeg', 'jpg', 'webp', 'gif', 'bmp'];
router.post('/v1/user/me/picture', async function (req, res) {
	// Check that the user exists
	const userQueryResult = await getUser(res.locals.user_id);
	if (!userQueryResult || userQueryResult.rows.length < 1) {
		return res.status(404).send({ error: 'User does not exists' });
	}

	// Get image from input
	if (req.files === undefined || req.files?.picture === undefined) {
		return res.status(400).send({ error: 'Picture missing' });
	}
	const picture = req.files.picture;
	if (Array.isArray(picture)) {
		return res.status(400).send({ error: 'Invalid picture' });
	}

	// Validate input
	if (picture.size > 2_000_000) {
		return res.status(400).send({ error: 'File too large, the limit is 2Mb' });
	}
	const extension = picture.name.split('.').pop() ?? '';
	if (!picture.mimetype.startsWith('image/') && allowedExtensions.indexOf(extension) < 0) {
		return res.status(400).send({ error: 'Invalid picture' });
	}

	// Also verify the magic mime type
	// try {
	// 	const type = await imageType(picture.data);
	// 	console.log(type);
	// } catch (error) {
	// 	console.error(error);
	// 	return res.status(400).send({ error: 'Invalid picture' });
	// }

	const path = `${CACHE_PATH}/${res.locals.user_id}.${extension}`;
	const absolutePath = resolve(path);
	try {
		// Save to disk
		await writeFile(absolutePath, picture.data);
		// Save to database
		await setUserExtension(res.locals.user_id, extension);
	} catch (error) {
		console.log(error);
		return res.status(500).send({ error: 'Failed to save image' });
	}

	return res.status(200).send({ status: 'ok' });
});

router.delete('/v1/user/me/picture', async function (req, res) {
	// Check that the user exists
	const userQueryResult = await getUser(res.locals.user_id);
	if (!userQueryResult || userQueryResult.rows.length < 1) {
		return res.status(404).send({ error: 'User does not exists' });
	}
	// Check if the user doesn't have a picture
	const user = userQueryResult.rows[0];
	if (!user.extension) {
		return res.status(200).send({ status: 'ok' });
	}

	const path = `${CACHE_PATH}/${res.locals.user_id}.${user.extension}`;
	const absolutePath = resolve(path);
	try {
		// Delete on disk
		await unlink(absolutePath);
		// Delete in database
		await setUserExtension(res.locals.user_id, null);
	} catch (error) {
		console.log(error);
		return res.status(500).send({ error: 'Failed to delete image' });
	}

	return res.status(200).send({ status: 'ok' });
});

router.get('/v1/user/:userId/picture', async function (req, res) {
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
