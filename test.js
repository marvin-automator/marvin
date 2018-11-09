time.cron.schedule({expression: "0-59/10 * * * * * *"}, (result) => {
    return http.request.send({
        url: "https://webhook.site/c192e8a8-6850-48ff-8eb7-bf7128b8d2f6?d=" + JSON.stringify(result),
        method: "GET"
    })
});