import path from 'path';
import { Router } from 'express';
import request from 'request';
import fs from 'fs';

const router = new Router();

router.get('/', async (req, res, next) => {
  var files = fs.readdirSync(path.resolve(__dirname, '../../qualifiers/'));

  res.send(JSON.stringify(files));
});

export default router;

