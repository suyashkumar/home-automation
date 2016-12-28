import { FETCH_TEMPS } from '../actions/index';
export default function(state = {}, action) {
	switch(action.type) {
		case FETCH_TEMPS:
			if (action.promiseLoading) {
				return Object.assign({}, state, {loading: true});
			} else {
				return Object.assign({}, state, {data: action.payload.data, loading: false});
			}
	}
	return state;
}
