import path from 'path';
import { Router } from 'express';
import request from 'request';
import fs from 'fs';

const router = new Router();

router.get('/', async (req, res, next) => {
  request({
    body: fs.readFileSync(path.resolve(__dirname, '../../qualifiers/problem_0.json')),
    method: 'POST',
    uri: 'http://localhost:8080/newgame'
  }, (err, resp, body) => {
    res.send(body);
  });
});

export default router;

