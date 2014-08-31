package main

import ()

var defaultTemplateIndexHtml = `<!DOCTYPE html>
<html>
	<head>
		<title>Kwikwik - {{.Title}}</title>
		
		<link type="text/css" rel="stylesheet" href="styles/default.css">
	</head>

	<body>
		<div class="page-title">{{.Title}}</div>

		<div class="menu">[ <a href="/edit/{{.Title}}">edit</a> ]</div>

		<div class="content-text">{{.Body}}</div>
	</body>
</html>
`
var defaultTemplateViewHtml = `<!DOCTYPE html>
<html>
	<head>
		<title>Kwikwik - {{.Title}}</title>

		<link type="text/css" rel="stylesheet" href="/styles/default.css">
		<link rel="icon" href="data:;base64,iVBORw0KGgo=">
	</head>

	<body>
		<div class="page-title">{{.Title}}</div>

		<div class="menu">[ <a href="/">index</a> ] [ <a href="/edit/{{.Title}}">edit</a> ]</div>

		<div class="content-text">{{.Body}}</div>
	</body>
</html>
`
var defaultTemplateEditHtml = `<!DOCTYPE html>
<html>
	<head>
		<title>Kwikwik - {{.Title}}</title>

		<link type="text/css" rel="stylesheet" href="/styles/default.css">
		<link rel="icon" href="data:;base64,iVBORw0KGgo=">
	</head>

	<body>
		<div class="page-title">Editing: {{.Title}}</div>

		<div class="menu">[ <a href="/">index</a> ]</div>

		<form action="/save/{{.Title}}" method="POST">
			<div><textarea name="body" rows="34" cols="120" class="editbox">{{printf "%s" .Body}}</textarea></div>
			<div><input type="submit" value="Save" class="button"></div>
		</form>
	</body>
</html>
`
var defaultTemplateErrorHtml = `<!DOCTYPE html>
<html>
	<head>
		<title>Kwikwik - Page not found</title>

		<link type="text/css" rel="stylesheet" href="/styles/default.css">
		<link rel="icon" href="data:;base64,iVBORw0KGgo=">
	</head>

	<body>
		<div class="page-title">Page not found</div>

		<div class="menu">[ <a href="/">index</a> ]</div>

		<div class="content-text">Page not found ( 404 ). </div>
	</body>
</html>
`
