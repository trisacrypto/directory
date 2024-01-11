import { Text, Button, HStack, Alert, AlertIcon, Box } from '@chakra-ui/react';

import { t, Trans } from '@lingui/macro';
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

const AttentionAlert = ({ severity, message, action, onClick }: AttentionAlertProps) => {
  if (severity === AttentionSeverity.INFO.toUpperCase()) {
    switch (action as AttentionAction) {
      case AttentionAction.START_REGISTRATION:
        return (
          <>
            <Alert bg="#D8EAF6" status={severity.toLowerCase()} borderRadius={'10px'}>
              <AlertIcon />
              <HStack justifyContent={'space-between'} w="100%">
                <Text> {message}</Text>
                <Button
                  onClick={onClick}
                  border={'1px solid white'}
                  width={142}
                  px={8}
                  as={'a'}
                  borderRadius={0}
                  color="#fff"
                  cursor="pointer"
                  bg="#000"
                  _hover={{ bg: '#000000D1' }}>
                  <Trans>Start</Trans>
                </Button>
              </HStack>
            </Alert>
          </>
        );

      case AttentionAction.COMPLETE_REGISTRATION:
        return (
          <Alert bg="#D8EAF6" status={severity.toLowerCase()} borderRadius={'10px'}>
            <AlertIcon />
            <HStack justifyContent={'space-between'} w="100%">
              <Text> {message}</Text>
              <Button
                onClick={onClick}
                width={142}
                border={'1px solid white'}
                px={8}
                as={'a'}
                borderRadius={0}
                background="#000"
                color="#fff"
                cursor="pointer"
                _hover={{ bg: '#000000D1' }}>
                <Trans>Complete</Trans>
              </Button>
            </HStack>
          </Alert>
        );
      case AttentionAction.SUBMIT_TESTNET:
      case AttentionAction.SUBMIT_MAINNET:
        return (
          <Alert bg="#D8EAF6" status={severity.toLowerCase()} borderRadius={'10px'}>
            <AlertIcon />
            <HStack justifyContent={'space-between'} w="100%">
              <Text minW={'70%'}> {t`${message}`}</Text>
              <Box>
                <Button
                  onClick={onClick}
                  border={'1px solid white'}
                  px={6}
                  as={'a'}
                  borderRadius={0}
                  background="#000"
                  color="#fff"
                  cursor="pointer"
                  _hover={{ bg: '#000000D1' }}>
                  {t`Submit`}
                </Button>
              </Box>
            </HStack>
          </Alert>
        );

      default:
        return (
          <Alert bg="#D8EAF6" status={AttentionSeverity.INFO} borderRadius={'10px'}>
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
