// Amazon Cognito creates a session which includes the id, access, and refresh tokens of an authenticated user.
const AmazonCognitoIdentity = require('amazon-cognito-identity-js')
const fs = require('fs')

const file = './http-client.private.env.json'
const testUserEmail = process.env.TEST_USER_EMAIL

const authenticationData = {
    Username : testUserEmail,
    Password : process.env.TEST_USER_PASSWORD,
};
const authenticationDetails = new AmazonCognitoIdentity.AuthenticationDetails(authenticationData);
const poolData = {
    UserPoolId : process.env.COGNITO_USER_POOL_ID,
    ClientId : process.env.COGNITO_CLIENT_ID
};
const userPool = new AmazonCognitoIdentity.CognitoUserPool(poolData);
const userData = {
    Username : testUserEmail,
    Pool : userPool
};
const cognitoUser = new AmazonCognitoIdentity.CognitoUser(userData);
cognitoUser.authenticateUser(authenticationDetails, {
    onSuccess: function (result) {
        const accessToken = result.getAccessToken().getJwtToken();

        fs.readFile(file, (err, data) => {
            if (err) {
                console.error(err)
                throw err
            }

            const json = JSON.parse(data.toString());
            json.dev.token = accessToken

            fs.writeFile(file, JSON.stringify(json, null, 2), (err) => {
                if (err) throw err
                console.log('Done')
            })
        })
    },

    onFailure: function(err) {
        console.log(err)
    },
});