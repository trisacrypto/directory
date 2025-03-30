# GDS Admin UI

The GDS Admin UI is a web application that allows TRISA reviewers to interact with the GDS Admin API and to manage the TRISA network. This application is built for both admin.trisa.directory and admin.testnet.directory.

## Getting Started

To develop the GDS Admin UI, create a `.env` file with the following:

```
REACT_APP_GDS_API_ENDPOINT=http://localhost:4434/v2
REACT_APP_GOOGLE_CLIENT_ID=
REACT_APP_GDS_IS_TESTNET=true
```

Note that you'll need to request the `$REACT_APP_GOOGLE_CLIENT_ID` configuration from one of the project leads who will whisper it to you.

In the root of the repository, run the GDS backend using docker compose:

```
$ docker compose -f ./containers/docker-compose.yaml --profile=gds build
$ docker compose -f ./containers/docker-compose.yaml --profile=gds up
```

Then in the project directory, install the dependencies defined in `package.json` and start the web server, which runs the application in development mode:

```
$ yarn
$ yarn start
```

Open [http://localhost:3000](http://localhost:3000) to view it in the browser. The page will reload if you make edits. You will also see any lint errors in the console.

To run the tests, use `yarn test`, which launches the test runner in the interactive watch mode. See the section about [running tests](https://facebook.github.io/create-react-app/docs/running-tests) for more information.

To build a production version, use `yarn build`, which builds the app for production to the `build` folder. It correctly bundles React in production mode and optimizes the build for the best performance. The build is minified and the filenames include the hashes. See the section about [deployment](https://facebook.github.io/create-react-app/docs/deployment) for more information.

Note that GDS-Admin-UI is deployed via a Docker container, which can be found at `containers/gds-admin-ui/Dockerfile` relative to the repository root.

## Available Scripts

- `yarn start`: Runs the app in development mode
- `yarn test`: Launches the test runner in interactive watch mode
- `yarn build`: Builds the app for production to the `build` folder
- `yarn eject`: Removes the build dependency (see notes below)
- `yarn flow`: Executes flow static type chcker
- `yarn format`: Executes prettier for format the code
- `yarn start:db`: runs a mock `json-server` with a mock database

## Other Notes

This project was bootstrapped with [Create React App](https://github.com/facebook/create-react-app). You can learn more in the [Create React App documentation](https://facebook.github.io/create-react-app/docs/getting-started). To learn React, check out the [React documentation](https://reactjs.org/).

### Node Sass

If you're getting an error that looks like:

> Error: Node Sass does not yet support your current environment: OS X 64-bit with Unsupported runtime

Then the [quickest fix](https://proustibat.medium.com/how-to-fix-error-node-sass-does-not-yet-support-your-current-environment-os-x-64-bit-with-c1b3298e4af0) is to run:

```
$ npm rebuild node-sass
```

Even though we're using `yarn` this seems to fix the problem locally on my machine without creating a `package-lock.json` file.

### `yarn eject`

**Note: this is a one-way operation. Once you `eject`, you can’t go back!**

If you aren’t satisfied with the build tool and configuration choices, you can `eject` at any time. This command will remove the single build dependency from your project.

Instead, it will copy all the configuration files and the transitive dependencies (webpack, Babel, ESLint, etc) right into your project so you have full control over them. All of the commands except `eject` will still work, but they will point to the copied scripts so you can tweak them. At this point you’re on your own.

You don’t have to ever use `eject`. The curated feature set is suitable for small and middle deployments, and you shouldn’t feel obligated to use this feature. However we understand that this tool wouldn’t be useful if you couldn’t customize it when you are ready for it.

### Guides

- [Code Splitting](https://facebook.github.io/create-react-app/docs/code-splitting)
- [Analyzing the Bundle Size](https://facebook.github.io/create-react-app/docs/analyzing-the-bundle-size)
- [Making a Progressive Web App](https://facebook.github.io/create-react-app/docs/making-a-progressive-web-app)
- [Advanced Configuration](https://facebook.github.io/create-react-app/docs/advanced-configuration)
- [Deployment](https://facebook.github.io/create-react-app/docs/deployment)
- [`yarn build` fails to minify](https://facebook.github.io/create-react-app/docs/troubleshooting#npm-run-build-fails-to-minify)