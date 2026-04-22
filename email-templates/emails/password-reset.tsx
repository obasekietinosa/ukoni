import {
  Body,
  Button,
  Container,
  Head,
  Html,
  Preview,
  Section,
  Text,
} from '@react-email/components';
import * as React from 'react';

interface PasswordResetEmailProps {
  resetLink?: string;
  userName?: string;
}

export const PasswordResetEmail = ({
  resetLink = '{{ .ResetLink }}',
  userName = '{{ .UserName }}',
}: PasswordResetEmailProps) => {
  return (
    <Html>
      <Head />
      <Preview>Reset your password</Preview>
      <Body style={main}>
        <Container style={container}>
          <Text style={title}>Password Reset</Text>
          <Text style={text}>
            Hi {userName},
          </Text>
          <Text style={text}>
            Someone recently requested a password change for your account. If this was you, you can set a new password here:
          </Text>
          <Section style={buttonContainer}>
            <Button style={button} href={resetLink}>
              Reset password
            </Button>
          </Section>
          <Text style={text}>
            If you don't want to change your password or didn't request this, just ignore and delete this message.
          </Text>
        </Container>
      </Body>
    </Html>
  );
};

export default PasswordResetEmail;

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

const buttonContainer = {
  padding: '24px 48px',
  textAlign: 'center' as const,
};

const button = {
  backgroundColor: '#007ee6',
  borderRadius: '4px',
  color: '#fff',
  fontSize: '16px',
  fontWeight: 'bold',
  textDecoration: 'none',
  textAlign: 'center' as const,
  display: 'block',
  width: '100%',
  padding: '12px',
};
