import React, { Component } from 'react';
import axios from '../axios-auth';
import {Button} from 'react-toolbox/lib/button';
import { Card, CardMedia, CardTitle, CardText, CardActions } from 'react-toolbox/lib/card';
import Input from 'react-toolbox/lib/input';
import constants from './constants';
import './login/login.css'; 

class Login extends Component {

	state = {
		email: "",
		password: "",
	}

	handleChange = (name, value) => {
		this.setState({[name]: value});
	}

	handleLogin = () => {
		axios.post(`${constants.serverUrl}/api/auth`, {
			email: this.state.email,
			password: this.state.password
		}).then(response => {
			if(response.data.success) {
				localStorage.setItem('jwtToken', response.data.token);
				window.open('/#/', "_self");
			} else {
				console.log('Error!');
			}
		});
	
	}
	
	render() {
		return ( 
			<div className="login">
				<Card className="login-container"> 
					<CardTitle title="Login" />	
					<div className="form">
						<Input 
							type="text" 
							label="Email" 
							name="email"
							value={this.state.email}
							onChange={this.handleChange.bind(this, 'email')}/>
						<Input 
							type="password" 
							label="Password" 
							name="password"
							value={this.state.password}
							onChange={this.handleChange.bind(this, 'password')}/>
						<Button className="button" onClick={this.handleLogin} raised primary>
							Login
						</Button>
						<Button className="button" raised>
							Register	
						</Button>
					</div>
				</Card>
			</div>
		);
	}
}

export default Login;
