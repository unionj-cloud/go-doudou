package registry

var indexTmpl = `<!DOCTYPE html>
<html lang="en">
<body>
<h1>Hello</h1>
<script>
  window.rows = {{ .Rows }}
</script>
</body>
</html>`