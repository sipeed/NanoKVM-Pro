# NanoKVM Frontend

This is NanoKVM web project. For more documentation, please refer to the [Wiki](https://wiki.sipeed.com/nanokvmpro).

## Structure

```shell
src
├── api                # backend api
├── assets             # static resources
├── components         # public components
├── i18n               # language resources
├── jotai              # Global jotai variables
├── lib                # util libs
├── pages              # web pages
│    ├── auth          # login and password
│    ├── desktop       # remote desktop
│    └── terminal      # web terminal
├── router.tsx         # routers
└── types              # types
```

## Build & Deploy

```shell
cd web
pnpm install

# After the compilation is complete, the `dist` folder will be generated
pnpm build

# Upload `dist` to `/kvmapp/server/` in NanoKVM-Pro
scp -r dist root@your-nanokvmpro-ip:/kvmapp/server/

# Replace the original `web` folder
rm -rf web
mv dist web

# Restart the service by executing
systemctl restart nanokvm
```
