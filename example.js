import http from "k6/http";
import { check, group } from "k6";

export let options = {
    stages: [
        { duration: "5s", target: 5 },
        { duration: "10s", target: 5 },
        { duration: "5s", target: 0 }
    ]
};

export default function() {
    // GET request
    group("GET", function() {
        let res = http.get("http://httpbin.org/get?verb=get");
        check(res, {
            "status is 200": (r) => r.status === 200,
            "is verb correct": (r) => r.json().args.verb === "get",
        });
    });


    // PUT request
    group("PUT", function() {
        let res = http.put("http://httpbin.org/put", JSON.stringify({ verb: "put" }), { headers: { "Content-Type": "application/json" }});
        check(res, {
            "status is 200": (r) => r.status === 200,
            "is verb correct": (r) => r.json().json.verb === "put",
        });
        console.log("my log")
    });

    // PATCH request
    group("PATCH", function() {
        let res = http.patch("http://httpbin.org/patch", JSON.stringify({ verb: "patch" }), { headers: { "Content-Type": "application/json" }});
        check(res, {
            "status is 200": (r) => r.status === 200,
            "is verb correct": (r) => r.json().json.verb === "patch",
        });
        console.error("my_error")
    });

    group("json-checker", function () {
        // Send a JSON encoded POST request
        let body = JSON.stringify({key: "value"});
        let res = http.post("http://httpbin.org/post", body, {headers: {"Content-Type": "application/json"}});

        // Use JSON.parse to deserialize the JSON (instead of using the r.json() method)
        let j = JSON.parse(res.body);

        // Verify response
        check(res, {
            "status is 200": (r) => r.status === 200,
            "is key correct": (r) => j.json.key === "value1",
            "Is stylesheet 4859 bytes?": (r) => r.body.length === 4859,
        });

        console.log(res.body.length)
    });

    // DELETE request
    group("DELETE-1", function() {
        let res = http.del("http://httpbin.org/delete?verb=delete");
        console.log("DELLL2")
        check(res, {
            "is verb correct": (r) => r.json().args.verb === "not-found-error-log-tets",
        });
    });

    // DELETE request
    group("DELETE-2", function() {
        let res = http.del("http://httpbin.org/1delete1?verb=delete");
        console.log("TESTTTTT1")
        check(res, {
            "status is 200": (r) => r.status === 200,
            "is verb correct": (r) => r.json().args.verb === "delete",
        });
    });
}