<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width,initial-scale=1.0">
    <base href="/">
    <link rel="icon" href="data:image/svg+xml,%3Csvg%20xmlns='http://www.w3.org/2000/svg'%20viewBox='0%200%2019%2019'%3E%3Ctext%20x='0'%20y='16'%3E🐻%3C/text%3E%3C/svg%3E" type="image/svg+xml" />
    <title>Gilk</title>
    <meta name="description" content="Gilk ~ Go per request query profiler 🐻">
    <link rel="stylesheet" href="/static/bulma.min.css">
    <link rel="stylesheet" href="/static/main.css">
    <script src="/static/alpine.min.js" defer></script>
</head>
<body>
    <div class="navbar">
        <span>Gilk</span>
        <span>Change theme</span>
    </div>
{{range $context := .}}
    <div>
        {{$context.Queries}}
    </div>
{{end}}
</body>
</html>