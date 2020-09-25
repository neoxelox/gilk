<!DOCTYPE html>
<html lang="en" class="has-navbar-fixed-top" style="background-color: #f5f5f5;">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width,initial-scale=1.0">
    <base href="/">
    <link rel="icon" href="/static/images/bear-192x192.png" type="image/png" sizes="192x192"/>
    <title>Gilk</title>
    <meta name="description" content="Gilk ~ Go per request query profiler 🐻">
    <!-- Twitter -->
    <meta name="twitter:card" content="summary">
    <meta name="twitter:title" content="Gilk 🐻">
    <meta name="twitter:description" content="Go per request query profiler 🔎">
    <meta name="twitter:image" content="https://raw.githubusercontent.com/Neoxelox/shortr/master/static/images/banner.png">
    <!-- Open Graph -->
    <meta property="og:type" content="summary">
    <meta property="og:site_name" content="Gilk">
    <meta property="og:title" content="Gilk 🐻">
    <meta property="og:description" content="Go per request query profiler 🔎">
    <meta property="og:image" content="https://raw.githubusercontent.com/Neoxelox/shortr/master/static/images/banner.png">
    <link rel="stylesheet" href="/static/styles/tomorrow.css">
    <script src="/static/scripts/highlight.pack.js"></script>
    <script>hljs.initHighlightingOnLoad();</script>
    <link rel="stylesheet" href="/static/styles/bulma.min.css">
    <link rel="stylesheet" href="/static/styles/main.css">
    <script src="/static/scripts/alpine.min.js" defer></script>
</head>
<body style="min-height: calc(100vh - 55px); background-color: #f5f5f5;">
    <nav class="navbar is-primary is-fixed-top" role="navigation" aria-label="main navigation" style="display: block;">
        <div class="navbar-brand">
            <a class="navbar-item" href="https://github.com/neoxelox/gilk">
                <img src="/static/images/bear-logo.png" width="32" height="32">
                <h1 class="title has-text-white" style="margin-left: 10px; margin-bottom: 2px;">Gilk</h1>
            </a>
            <div class="navbar-end" style="margin-left: auto;">
                <div class="navbar-item" style="height: 100%;">
                    <a href="https://github.com/Neoxelox">
                        <div class="tags has-addons" style="margin-bottom: 0px; margin-right: 10px;">
                            <span class="tag" style="margin-bottom: 0px;">Author</span>
                            <span class="tag is-primary is-light" style="margin-bottom: 0px;">@Neoxelox</span>
                        </div>
                    </a>
                    <div class="tags has-addons">
                        <span class="tag">Version</span>
                        <span class="tag is-primary is-light">0.0.1</span>
                    </div>
                </div>
            </div>
        </div>
    </nav>

    <div class="columns is-multiline is-mobile" style="margin: 15px 15px 0px 15px;">
    {{range $cindex, $context := .}}
    {{if $context.HasFinished}}
        <div class="column">
            <div x-data="{ tab: 'overall' }" class="card" style="border-radius: 6px; min-width: 375px;">
            <header class="card-header">
                <p class="card-header-title">
                {{$context.Path}}
                </p>
                <span class="card-header-icon" aria-label="more options" style="cursor: default;">
                <span class="tag is-{{$context.MethodColor}}">{{$context.Method}}</span>
                </span>
            </header>
            <div class="card-content" :class="{ 'queries': tab === 'queries' }">
                <div x-show="tab === 'overall'" class="content">
                    <center>
                        <h3 class="title is-3" style="margin: 0; margin-top: 1rem; margin-bottom: 0.5rem; line-height: 0;"><span class="has-text-{{$context.ContextColor}}">{{$context.ContextDuration}}</span> overall</h3> <br> 
                        <hr style="margin: 0; margin-bottom: 1rem; line-height: 0;">
                        <h3 class="title is-3" style="margin: 0;"><span class="has-text-{{$context.QueriesColor}}">{{$context.QueriesDuration}}</span> on queries</h3> <br>
                        <h4 class="subtitle is-4" style="margin: 0; margin-bottom: 1rem; line-height: 0;"><span class="has-text-{{$context.LenQueriesColor}}">{{len $context.Queries}}</span> queries</h4>
                    </center>
                </div>
                <div x-show="tab === 'queries'" class="content">
                    <div class="accordion">
                    {{range $qindex, $query := $context.Queries}}
                        <div class="accordion-tab">
                            <input type="checkbox" id="{{$cindex}}-{{$qindex}}">
                            <label class="accordion-tab-label" for="{{$cindex}}-{{$qindex}}">
                                <span>
                                    {{$query.CallerFunc}}
                                    <span class="tag is-{{$query.Color}} is-light" style="margin-left: 10px; transform: translateY(-1px);">{{$query.Duration}}</span>
                                </span>
                            </label>
                            <div class="accordion-tab-content">
                                <div class="snippet">
                                    <pre><code class="sql">{{$query.Format}}
                                        <small>--{{$query.CallerFile}}:{{$query.CallerLine}}</small>
                                    </code></pre>
                                </div>
                            </div>
                        </div>
                    {{end}}
                    </div>
                </div>
            </div>
            <footer class="card-footer">
                <div class="tabs has-text-primary">
                    <ul>
                      <li :class="{ 'is-active': tab === 'overall' }" @click="tab = 'overall'"><a>Overall</a></li>
                      <li :class="{ 'is-active': tab === 'queries' }" @click="tab = 'queries'"><a>Queries</a></li>
                    </ul>
                  </div>
            </footer>
            </div>
        </div>
    {{end}}
    {{end}}
    </div>
</body>
</html>