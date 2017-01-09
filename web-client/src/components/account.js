import React, { Component } from 'react';
import axios from '../axios-auth';
import { Card, CardMedia, CardTitle, CardText, CardActions } from 'react-toolbox/lib/card';
import constants from './constants';
import copy from 'copy-to-clipboard';
import { List, ListItem, ListSubHeader, ListDivider, ListCheckbox } from 'react-toolbox/lib/list';
import './account/account.css';
import './common.css';

class Account extends Component {
	constructor(props) {
		super(props);
		this.state = {
			email: "",
			apiKey: "",
		}
		this.getAccountInformation();
	}

	getAccountInformation = () => {
		axios.get(`${constants.serverUrl}/api/me`, {auth_me: true}).then(response => {
			if (response.data.success) {
				this.setState({email: response.data.email, apiKey: response.data.key});
			}
			console.log(response);
		}); 
	}

	render() {
		return (
			<div className="container">
				<Card className="content">
					<CardTitle title="Account" />
					<List selectable>
						<ListItem
							caption="Email"
							legend={this.state.email}/>
						<ListItem
							caption="API Key (Click to copy)"
							legend={this.state.apiKey}
							onClick={() => copy(this.state.apiKey)}/>
					</List>
				</Card>
			</div>
		);
	}

}

export default Account;
