import React from 'react';
import ReactDOM from 'react-dom';
import Account from './components/account';
import { Provider } from 'react-redux';
import { createStore, applyMiddleware } from 'redux';
import ReduxPromise from 'redux-promise';
import App from 'components/app';
import './index.css';
import reducers from './reducers';
import { Router, Route, Link, IndexRoute, hashHistory} from 'react-router'
import ReduxLoadingPromise from './middleware/redux-loading-promise.js';
import Interact from './components/interact';
import Login from './components/login';
import Home from './components/home';

const createStoreWithMiddleware = applyMiddleware(ReduxLoadingPromise)(createStore);

ReactDOM.render(
	<Provider store={createStoreWithMiddleware(reducers)}>
		<Router history={hashHistory}>
			<Route path="/" component={App}>
				<IndexRoute component={Home} />
				<Route path="/login" component={Login} />
				<Route path="/interact" component={Interact} />
				<Route path="/account" component={Account} />
			</Route>
		</Router> 
	</Provider>, 
 	document.getElementById('root')
);
