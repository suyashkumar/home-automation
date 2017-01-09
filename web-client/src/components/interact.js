import React, { Component } from 'react';
import axios from '../axios-auth';
import {Button} from 'react-toolbox/lib/button';
import { Card, CardMedia, CardTitle, CardText, CardActions } from 'react-toolbox/lib/card';
import Input from 'react-toolbox/lib/input';
import constants from './constants';
import './interact/interact.css';
import './common.css';

class Interact extends Component {

	state = {
		deviceName: '',
		functionName: '',
		latestResponse: '',
	}

	handleChange = (name, value) => {
		this.setState({[name]: value});
	}

	callFunction = () => { 
		axios.get(`${constants.serverUrl}/api/send/${this.state.deviceName}/${this.state.functionName}`, {auth_me: true}).then(response => {
			if (response.data.success) {
				this.setState({latestResponse: response.data.data});
			}
			console.log(response);
		});
	}

	//TODO (suyashkumar) improve this logic
	handleKeyPress = e => {
		if (e.key === 'Enter') {
			this.callFunction();
		}
	}

	render() {
		return (
			<div className="container">
				<Card className="content">
					<CardTitle title="Interact"/>
					<div className="">
						<p>
							Enter a device name and a function name below and hit "Go"
						</p>
						<Input
							type="text"
							label="Device Name"
							name="deviceName" 
							className="form-input"
							value={this.state.deviceName}
							onChange={this.handleChange.bind(this, 'deviceName')}/>
						<Input
							type="text"
							label="Function Name"
							name="functionName" 
							className="form-input"
							value={this.state.functionName}
							onKeyPress={this.handleKeyPress}
							onChange={this.handleChange.bind(this, 'functionName')}/>
						{
							this.state.latestResponse &&
							<div className="response-area">
								<h4>Response:</h4>
								<div className="response-content">
									{this.state.latestResponse}
								</div>
							</div>
						}
						<Button 
						className="button" 
						onClick={this.callFunction} raised primary>
							Go!
						</Button>
					</div>
				</Card>
			</div>
		);
	}

}

export default Interact;
