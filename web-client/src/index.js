import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import { createStore, applyMiddleware } from 'redux';
import ReduxPromise from 'redux-promise';
import App from 'components/app';
import './index.css';
import reducers from './reducers';
import ReduxLoadingPromise from './middleware/redux-loading-promise.js';

const createStoreWithMiddleware = applyMiddleware(ReduxLoadingPromise)(createStore);

ReactDOM.render(
	<Provider store={createStoreWithMiddleware(reducers)}>
		<App />
	</Provider>, 
 	document.getElementById('root')
);
