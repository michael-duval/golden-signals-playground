import http from 'k6/http';
export const options = {
    scenarios: {
        // arbitrary name of scenario
        average_load: {
          executor: 'ramping-vus',
          stages: [
            // ramp up to average load of 20 virtual users
            { duration: '30s', target: 20 },
            // maintain load
            { duration: '2m', target: 20 },
            // ramp down to zero
            { duration: '15s', target: 0 },
          ],
        },
    }
}

export default function () {
    const latency = 50;
    const failPct = 0.60;

    const url = `http://localhost:8080/work?latency_ms=${latency}&fail=${failPct}`;
    http.get(url)
}