# GDS UI

The GDS UI is a web application that allows users to interact with the GDS and register to join the TRISA network. This application is built for both trisa.directory and testnet.directory.

## Getting Started

To develop the GDS UI, create a `.env` file with the following:

```
REACT_APP_GDS_API_ENDPOINT=http://localhost:8080
REACT_APP_GDS_IS_TESTNET=true
REACT_APP_ANALYTICS_ID=[GOOGLE ANALYTICS ID]
```

Note that the `$REACT_APP_ANALYTICS_ID` is not required.

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

Note that GDS-UI is deployed via a Docker container, which can be found at `containers/gds-ui/Dockerfile` relative to the repository root.

## Available Scripts

- `yarn start`: Runs the app in development mode
- `yarn test`: Launches the test runner in interactive watch mode
- `yarn build`: Builds the app for production to the `build` folder
- `yarn eject`: Removes the build dependency (see notes below)
- `yarn protos`: Regenerates the protocol buffers (see notes below)
- `yarn extract`: Extracts translated text using `lingui extract`
- `yarn extract-c`: Extracts translated text using `lingui extract --clean`
- `yarn compile`: Compiles translated .po files using `lingui compile`

## Protocol Buffers

The GDS-UI uses [protocol buffers](https://developers.google.com/protocol-buffers/docs/reference/javascript-generated) and [grpc-web](https://github.com/grpc/grpc-web) for its backend communication. As a result, the protocol buffers in `protos` and from the `trisacrypto/trisa` repository need to be available and generated to the latest version.

1. Install [protoc](https://grpc.io/docs/protoc-installation/) and [protoc-gen-grpc-web](https://github.com/grpc/grpc-web/releases/tag/1.3.0). The simplest way to do this is with Homebrew:

   ```
   $ brew install protoc
   $ brew install protoc-gen-grpc-web
   ```

2. Ensure that the TRISA repositories are cloned in the same workspace.

   ```
   $ cd ~/my/workspace
   $ git clone git@github.com:trisacrypto/directory.git
   $ git clone git@github.com:trisacrypto/trisa.git
   $ cd directory/web/gds-ui
   ```

3. Run the `yarn protos` command to generate the protocol buffers.

The `yarn protos` command runs the `bin/generate.sh` script to execute `protoc` directly. The generated files are placed in `src/api`.

## i18n steps

The GDS-UI is translated into several languages using [lingui](https://lingui.js.org/index.html). The basic workflow is as follows:

- Wrap all text messages that need to be translated with `<Trans> Message to be translated</Trans>`. Also add `import { Trans } from "@lingui/macro"` to each file containing wrapped messages. For text inside a select menu, please use ``i18n._(t`text`)`` and also remember to add `import { i18n } from "@lingui/core"; import { t } from "@lingui/macro";` for each file.

- Run `yarn extract` to extract all messages. Or run `yarn extract-c` to extract all messages while also removing the translations that are no longer in the source file.

- Send `.po` files to translators. The files can be found the the `"path": "src/locales/{locale}/messages"` path specified above.

- Update the `.po` files with translations received from the translators and run `yarn compile` to incorporate the translations.

### Installation Notes

The following is how we installed and setup i18n in this project.

- Install `lingui/cli`, `@lingui/macro`, and `@lingui/react1

   ```
   $ yarn add -dev @lingui/cli @babel/core
   $ yarn add -dev @lingui/macro babel-plugin-macros
   $ yarn add @lingui/react
   ```

- Create `.linguirc` with LinguiJS configuration, and place it next to `package.json`. Modify `locales` to contain all of your locales. Replace `src` with the target folder for internationalization.

   ```
   {
      "locales": ["en", "cs"],
      "catalogs": [{
         "path": "src/locales/{locale}/messages",
         "include": ["src"]
      }],
      "format": "po"
   }
   ```

- Add the following to `package.json`

   ```
   {
      "scripts": {
         "extract": "lingui extract",
         "compile": "lingui compile",
      }
   }
   ```

### Plural usage

In general, there are 6 plural forms (based on [CLDR](http://cldr.unicode.org/index/cldr-spec/plural-rules)]):

- zero

- one (singular)

- two (dual)

- few (paucal)

- many (also used for fractions if they have a separate class)

- other (required—general plural form—also used if the language only has a single form)

To use the plural forms in user messaging, e.g., [n] books, we can wrap our message as:

```
i18n.plural({
  value: numBooks,
  one: "# book",
  other: "# books"
})
```

When extracted by lingui command, the message is formatted as `{numBooks, plural, one {# book} other {# books}}` and the translators will need to follow the [plural forms](https://unicode-org.github.io/cldr-staging/charts/latest/supplemental/language_plural_rules.html) in their target language. E.g., if translating this message into Czech, it should become `{numBooks, plural, one {# kniha} few {# knihy} other {# knih}}`.

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