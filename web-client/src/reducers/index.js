import { combineReducers } from 'redux';
import TempReducer from './TempsReducer';
import LocationsReducer from './locations-reducer';

const rootReducer = combineReducers({
	temps: TempReducer,
	locationData: LocationsReducer
});

export default rootReducer;
