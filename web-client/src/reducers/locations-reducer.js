import { FETCH_LOCATIONS, SELECT_LOCATION } from '../actions/index';

export default (state={}, action) => {
	switch(action.type) {
		case FETCH_LOCATIONS:
			return Object.assign({}, state, {locations: action.payload.data});
		case SELECT_LOCATION:
			return Object.assign({}, state, {currentLocation: action.currentLocation});
	}
	return state;
}
