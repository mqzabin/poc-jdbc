import { check } from 'k6';
import http from 'k6/http';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import { randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

export const options = {
    discardResponseBodies: true,
};

const config = {
    host: "http://app-server:3000",
    headers: {
        'Content-Type': 'application/json',
    },
}

export function setup() {

}

export default function (customers) {
    const { host, headers } = config;

    const url = `${host}/create_event`;

    const params = { headers: headers };

    const payload = JSON.stringify({
        user_id: uuidv4(),
        product_id: uuidv4(),
        investment_id: uuidv4(),
        amount: randomIntBetween(1, 10000)
    });


    const res = http.post(url, payload, params);
    check(res, { 'is status 201': r => r.status === 201});
}