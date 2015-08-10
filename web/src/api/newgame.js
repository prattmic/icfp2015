import path from 'path';
import { Router } from 'express';
import request from 'request';
import fs from 'fs';

const router = new Router();

router.get('/', async (req, res, next) => {
  let qualifierName = req.query.qualifier || 'problem_0.json';

  let data = {
    Problem: JSON.parse(
      fs.readFileSync(path.resolve(__dirname, '../../qualifiers/' + qualifierName))
    ),
    AI: req.query.ai,
    Repeater: req.query.repeater
  };

  console.log('sending request');
  request({
    body: JSON.stringify(data),
    method: 'POST',
    uri: 'http://localhost:8080/newgame'
  }, (err, resp, body) => {
    console.log('received response', body && body.length);
    res.send(body);
  });
});

export default router;

