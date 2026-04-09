This is a small project to practice the http in go. 
in this project I built the backend of social media platform similar to twitter with limited featuers 

Features
--------
1. users can create a profile - "POST /api/users"
2. login - "POST /api/login"
3. create chirps - "POST /api/chirps"
4. see their metrics - "GET /admin/metrics"
5. refresh and revoke their login tokens - "POST /api/refresh" & "POST /api/revoke"
6. see the health of the site - "GET /api/healthz"
7. delete chirps - "DELETE /api/chirps/{chirpID}"
8. Upgrade to a premium account - "POST /api/polka/webhooks"
9. reset - reset the chirps
