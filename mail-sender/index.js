const { Resend } = require('resend');
require('dotenv').config();
var amqp = require('amqplib/callback_api');


const resend = new Resend(process.env.RESEND_API_KEY);
const QUEUE_NAME = 'email_queue';

amqp.connect('amqp://localhost', function(error0, connection) {
  if(error0){
    throw error0;
  }
  connection.createChannel(function(error1, channel){
    if(error1){
      throw error1;
    }
    channel.assertQueue(QUEUE_NAME, {
      durable: true 
    });
    channel.prefetch(1);
    console.log("Successful, waiting for signal to send an email") // for testing
    channel.consume(QUEUE_NAME, async function(data) {
      try {
      await sendEmail(data);
      channel.ack(data);
      } catch(error) {
        console.error('Error sending email:', error);
        channel.nack(data,false,true);
      }
      }, {
          noAck: false 
        });
      });
    });

async function sendEmail(data) {
    //here we will use values for our email from the producer(json/etc...)
    //Example:
    //from = data.from, and then from: from
   const verificationLink = `http://localhost:5000/verify?=${data.token}`;
  //const verificationLink = `http://localhost:5000/verify?=hello`;
    const email = await resend.emails.send({
      from: 'Acme <onboarding@resend.dev>',
      to: ['podbodniykirill@gmail.com'],
      subject: 'Email verification',
      html: `<h2>Email Verification</h2>
    <p>Click the button below to verify your email:</p>
<a href="${verificationLink}" style="padding:10px 20px; background:black; color:white;text-decoration:none;">
Verify email</a>`
    });

    console.log("Email sent:" + email);
};

