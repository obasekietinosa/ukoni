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

interface InvitationEmailProps {
  inviteLink?: string;
}

export const InvitationEmail = ({
  inviteLink = '{{ .InviteLink }}',
}: InvitationEmailProps) => {
  return (
    <Html>
      <Head />
      <Preview>You have been invited to join an inventory</Preview>
      <Body style={main}>
        <Container style={container}>
          <Text style={title}>Inventory Invitation</Text>
          <Text style={text}>
            You have been invited to join an inventory on Ukoni. Click the following button to accept the invitation:
          </Text>
          <Section style={buttonContainer}>
            <Button style={button} href={inviteLink}>
              Accept Invitation
            </Button>
          </Section>
        </Container>
      </Body>
    </Html>
  );
};

export default InvitationEmail;

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
