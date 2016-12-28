import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { fetchTemps, fetchLocations } from '../actions/index';
import { Button } from 'react-toolbox/lib/button';
import {Tab, Tabs} from 'react-toolbox';
import {LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend} from 'recharts';

class Test extends Component {
	
	render() {
		console.log(this.props.temps, this.props.locations);
	return (
		<div style={{textAlign: 'center', marginTop: '20px'}}> 
			<Button onClick={()=>this.props.fetchTemps('ADPLKenyaC6304')} raised >
				Test AJAX
			</Button>
			<Button onClick={this.props.fetchLocations} raised>
				Test Locations AJAX
			</Button>
			<LineChart data={this.props.temps[0] && this.props.temps[0].filter(val=>val.probeid=="HXCI")} width={600} height={300}>
				<XAxis dataKey="time"/> 
				<Tooltip/>
				<YAxis />
				<Line dataKey="temp" />
			</LineChart>

		</div>
	);
	}
}

const mapDispatchToProps = dispatch => {
	return bindActionCreators({ fetchTemps, fetchLocations }, dispatch);
}

const mapStateToProps = state => {
	return { 
		temps: state.temps,
		locations: state.locationData.locations
	}
}

export default connect(mapStateToProps, mapDispatchToProps)(Test);
