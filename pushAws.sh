## zip it up
zip  lambdaTest.zip lambdaTest
## upload it to aws
aws lambda update-function-code --function-name TestgoLambda --zip-file fileb://lambdaTest.zip --output text
## create a new version
aws lambda publish-version --function-name TestgoLambda --output text
