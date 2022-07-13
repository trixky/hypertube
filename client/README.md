# Install node/npm/npx

https://github.com/nodesource/distributions#debinstall

> use version "Node.js v16.x"

```
# Using Ubuntu
curl -fsSL https://deb.nodesource.com/setup_16.x | sudo -E bash -
sudo apt-get install -y nodejs

# Using Debian, as root
curl -fsSL https://deb.nodesource.com/setup_16.x | bash -
apt-get install -y nodejs
```

# Setup the svelte project

https://kit.svelte.dev/docs/introduction

svelte init option:

- typescript
- eslint
- prettier

```
npm init svelte client
cd client
npm install
npm run dev
```

# Setup tailwind

https://tailwindcss.com/docs/guides/sveltekit
