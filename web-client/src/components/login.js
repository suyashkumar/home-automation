import React, { Component } from 'react';
import {Button} from 'react-toolbox/lib/button';
import { Card, CardMedia, CardTitle, CardText, CardActions } from 'react-toolbox/lib/card';
import Input from 'react-toolbox/lib/input';
import './login/login.css';

class Login extends Component {

	state = {
		email: "",
		password: "",
	}

	handleChange = (name, value) => {
		this.setState({[name]: value});
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
						<Button className="button" raised primary>
							Login
						</Button>
					</div>
				</Card>
			</div>
		);
	}
}

export default Login;
