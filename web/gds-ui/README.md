# Getting Started with Create React App

This project was bootstrapped with [Create React App](https://github.com/facebook/create-react-app).

## Available Scripts

In the project directory, you can run:

### `yarn start`

Runs the app in the development mode.\
Open [http://localhost:3000](http://localhost:3000) to view it in the browser.

The page will reload if you make edits.\
You will also see any lint errors in the console.

### `yarn test`

Launches the test runner in the interactive watch mode.\
See the section about [running tests](https://facebook.github.io/create-react-app/docs/running-tests) for more information.

### `yarn build`

Builds the app for production to the `build` folder.\
It correctly bundles React in production mode and optimizes the build for the best performance.

The build is minified and the filenames include the hashes.\
Your app is ready to be deployed!

See the section about [deployment](https://facebook.github.io/create-react-app/docs/deployment) for more information.

### `yarn eject`

**Note: this is a one-way operation. Once you `eject`, you can’t go back!**

If you aren’t satisfied with the build tool and configuration choices, you can `eject` at any time. This command will remove the single build dependency from your project.

Instead, it will copy all the configuration files and the transitive dependencies (webpack, Babel, ESLint, etc) right into your project so you have full control over them. All of the commands except `eject` will still work, but they will point to the copied scripts so you can tweak them. At this point you’re on your own.

You don’t have to ever use `eject`. The curated feature set is suitable for small and middle deployments, and you shouldn’t feel obligated to use this feature. However we understand that this tool wouldn’t be useful if you couldn’t customize it when you are ready for it.

## Learn More

You can learn more in the [Create React App documentation](https://facebook.github.io/create-react-app/docs/getting-started).

To learn React, check out the [React documentation](https://reactjs.org/).

### Code Splitting

This section has moved here: [https://facebook.github.io/create-react-app/docs/code-splitting](https://facebook.github.io/create-react-app/docs/code-splitting)

### Analyzing the Bundle Size

This section has moved here: [https://facebook.github.io/create-react-app/docs/analyzing-the-bundle-size](https://facebook.github.io/create-react-app/docs/analyzing-the-bundle-size)

### Making a Progressive Web App

This section has moved here: [https://facebook.github.io/create-react-app/docs/making-a-progressive-web-app](https://facebook.github.io/create-react-app/docs/making-a-progressive-web-app)

### Advanced Configuration

This section has moved here: [https://facebook.github.io/create-react-app/docs/advanced-configuration](https://facebook.github.io/create-react-app/docs/advanced-configuration)

### Deployment

This section has moved here: [https://facebook.github.io/create-react-app/docs/deployment](https://facebook.github.io/create-react-app/docs/deployment)

### `yarn build` fails to minify

This section has moved here: [https://facebook.github.io/create-react-app/docs/troubleshooting#npm-run-build-fails-to-minify](https://facebook.github.io/create-react-app/docs/troubleshooting#npm-run-build-fails-to-minify)

## i8n steps

### Install
- Install `lingui/cli`, `@lingui/macro`, and `@lingui/react1
```
npm install --save-dev @lingui/cli @babel/core
npm install --save-dev @lingui/macro babel-plugin-macros
npm install --save @lingui/react
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

- Wrap all text messages that need to be translated with `<Trans> Message to be translated</Trans>`. Also add `import { Trans } from "@lingui/macro"` to each files containing wrapped messages.  

- Run `npm run extract` to extract all messages. Or run `npm run extract-c` to extract all messages while also removing the translations that are no longer in the source file. 

- Send `.po` files to translators. The files can be found the the `"path": "src/locales/{locale}/messages"` path specified above.

- Update the `.po` files with translations received from the translators and run `npm run compile` to incorporate the translations. 

