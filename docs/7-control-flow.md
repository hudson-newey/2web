# Control Flow

## If Conditions

```html
<script compiled>
  $condition = false;
</script>

<button @click="$condition = !$condition"></button>

@if ($condition) {
  <div>
    <h1>Discoverable content</h1>
  </div>
}
```

[**Next**](./8-code-splitting.md) (Code Splitting)
