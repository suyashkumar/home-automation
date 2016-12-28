const isPromise = possiblePromise => {
	return possiblePromise && typeof possiblePromise.then === 'function';
};

export default ({ dispatch }) => {
	return next => action => {
		console.log(isPromise(action.payload));
		console.log(action);
		if (isPromise(action.payload) && !action.promiseLoading) {
			// Dispatch the action immediately with a loading parameter
			dispatch({ ...action, promiseLoading: true});
			console.log('dispatched first');
			// Schedule dispatch for completed request with resultant payload
			action.payload.then(
				result => {dispatch({...action, payload: result, promiseLoading: false}); console.log('second')},
				error => {
					dispatch({...action, payload: error, error: true});
					return Promise.reject(error);
				}
			);
		} else {
			next(action);
		}

	};

}
