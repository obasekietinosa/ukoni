# Ukoni Email Templates

This is a [React Email](https://react.email) project used to design and build email templates for the Ukoni backend API.

## Workflow

We use React to design these emails so that we can easily preview them, type check them, and utilize common React patterns. However, the final emails must be shipped as raw HTML strings embedded in the Go backend.

1. **Design and Development**: Write or modify React components in the `emails/` directory.
2. **Local Preview**: Run `npm run dev` to start a local preview server at `localhost:3000`. You can test variables by supplying default values in the component props.
3. **Build**: Run `npm run build`. This script will:
   - Export the templates as pure HTML into the `out/` folder.
   - Automatically copy the generated HTML into `../api/internal/mailer/templates/`.

Once copied over, the backend will embed these `.html` files in the binary using `go:embed` and populate them with standard `html/template` or `text/template` Go syntax.

### Using Go Variables in React Email

Because the exported files act as standard Go templates (`html/template`), you can use Go text variables directly inside strings in React like so:

```tsx
export const PasswordResetEmail = ({
  resetLink = '{{ .ResetLink }}',
  userName = '{{ .UserName }}',
}: PasswordResetEmailProps) => {
  return (
    <Html>
      <Text>Hi {userName},</Text>
      <Button href={resetLink}>Reset password</Button>
    </Html>
  );
};
```

React Email will sometimes wrap text strings in `<!-- -->` (e.g. `<!-- -->{{ .UserName }}<!-- -->`) during compilation, which is perfectly safe and ignores perfectly when compiled by Go's `html/template`.

## Important Note

The `out/` directory here is strictly ignored by `.gitignore` to avoid duplication. The generated HTML inside `api/internal/mailer/templates/` **must** be committed to Git.
