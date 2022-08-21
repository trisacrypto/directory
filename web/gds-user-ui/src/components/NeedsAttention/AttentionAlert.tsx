import {
  Box,
  Text,
  Stack,
  Button,
  HStack,
  Alert,
  AlertIcon,
  AlertTitle,
  AlertDescription
} from '@chakra-ui/react';

import { NavLink } from 'react-router-dom';
import * as Sentry from '@sentry/react';
import { Trans } from '@lingui/react';
import { NeedsAttentionProps } from '.';
const enum AttentionSeverity {
  INFO = 'info',
  WARNING = 'warning',
  ERROR = 'error'
}
const enum AttentionAction {
  NO_ACTION = 'NO_ACTION',
  START_REGISTRATION = 'START_REGISTRATION',
  COMPLETE_REGISTRATION = 'COMPLETE_REGISTRATION',
  SUBMIT_TESTNET = 'SUBMIT_TESTNET',
  SUBMIT_MAINNET = 'SUBMIT_MAINNET',
  VERIFY_EMAILS = 'VERIFY_EMAILS',
  RENEW_CERTIFICATE = 'RENEW_CERTIFICATE',
  CONTACT_SUPPORT = 'CONTACT_SUPPORT'
}

export type AttentionResponseType = {
  message: string;
  severity: any;
  action: string;
};

type AttentionAlertProps = Partial<AttentionResponseType & NeedsAttentionProps>;

const AttentionAlert = ({
  severity,
  message,
  action,
  onClick,
  buttonText
}: AttentionAlertProps) => {
  if (severity === AttentionSeverity.INFO.toUpperCase()) {
    switch (action as AttentionAction) {
      case AttentionAction.START_REGISTRATION:
        return (
          <>
            <Alert status={severity.toLowerCase()} borderRadius={'10px'}>
              <AlertIcon />
              <HStack justifyContent={'space-between'}>
                <Text> {message}</Text>
                <Button
                  onClick={onClick}
                  border={'1px solid white'}
                  width={142}
                  px={8}
                  as={'a'}
                  borderRadius={0}
                  background="transparent"
                  color="#fff"
                  cursor="pointer"
                  _active={{ background: '#000' }}
                  _hover={{ background: '#111', color: 'white' }}>
                  Start
                </Button>
              </HStack>
            </Alert>{' '}
          </>
        );
      case AttentionAction.SUBMIT_TESTNET || AttentionAction.SUBMIT_MAINNET:
        return (
          <>
            <Alert status={severity.toLowerCase()} borderRadius={'10px'}>
              <AlertIcon />
              <HStack justifyContent={'space-between'}>
                <Text> {message}</Text>
                <Button
                  onClick={onClick}
                  border={'1px solid white'}
                  width={142}
                  px={8}
                  as={'a'}
                  borderRadius={0}
                  background="transparent"
                  color="#fff"
                  cursor="pointer"
                  _active={{ background: '#000' }}
                  _hover={{ background: '#111', color: 'white' }}>
                  {buttonText}
                </Button>
              </HStack>
            </Alert>{' '}
          </>
        );
      case AttentionAction.COMPLETE_REGISTRATION:
        return (
          <>
            <Alert status={severity.toLowerCase()} borderRadius={'10px'}>
              <AlertIcon />
              {message}
            </Alert>
            <Alert status={severity.toLowerCase()} borderRadius={'10px'}>
              <AlertIcon />
              <HStack justifyContent={'space-between'}>
                <Text> {message}</Text>
                <Button
                  onClick={onClick}
                  width={142}
                  border={'1px solid white'}
                  px={8}
                  as={'a'}
                  borderRadius={0}
                  background="transparent"
                  color="#fff"
                  cursor="pointer"
                  _active={{ background: '#000' }}
                  _hover={{ background: '#111', color: 'white' }}>
                  Complete
                </Button>
              </HStack>
            </Alert>{' '}
          </>
        );

      default:
        return (
          <Alert status={severity.toLowerCase()} borderRadius={'10px'}>
            <AlertIcon />
            {message}
          </Alert>
        );
    }
  } else {
    return (
      <Alert status={severity.toLowerCase()} borderRadius={'10px'}>
        <AlertIcon />
        {message}
      </Alert>
    );
  }
};

export default AttentionAlert;
