import {
  Body,
  Container,
  Head,
  Html,
  Preview,
  Text,
} from '@react-email/components';
import * as React from 'react';

interface WelcomeEmailProps {
  userName?: string;
}

export const WelcomeEmail = ({
  userName = '{{ .UserName }}',
}: WelcomeEmailProps) => {
  return (
    <Html>
      <Head />
      <Preview>Welcome to Ukoni</Preview>
      <Body style={main}>
        <Container style={container}>
          <Text style={title}>Welcome to Ukoni!</Text>
          <Text style={text}>
            Hi {userName},
          </Text>
          <Text style={text}>
            We're excited to have you on board. Enjoy using Ukoni!
          </Text>
        </Container>
      </Body>
    </Html>
  );
};

export default WelcomeEmail;

const main = {
  backgroundColor: '#f6f9fc',
  fontFamily:
    '-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,"Helvetica Neue",Ubuntu,sans-serif',
};

const container = {
  backgroundColor: '#ffffff',
  margin: '0 auto',
  padding: '20px 0 48px',
  marginBottom: '64px',
};

const title = {
  fontSize: '24px',
  lineHeight: '1.25',
  fontWeight: '700',
  color: '#333',
  padding: '0 48px',
};

const text = {
  color: '#333',
  fontSize: '16px',
  lineHeight: '24px',
  padding: '0 48px',
};
