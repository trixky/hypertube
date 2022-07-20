import express from 'express';
import get_position from '../handlers/get_position.handler';
import save_position from '../handlers/save_position.handler';

const router = express.Router();

router.get('/:torrent_id', get_position);
router.post('/:torrent_id', save_position);

export default router;
