import React, { FC } from 'react';
import {
  Tr,
  Td,
  Flex,
  Text,
  Icon,
  IconButton,
  Box,
  Menu,
  Button,
  useColorModeValue,
  MenuItem
} from '@chakra-ui/react';

import { IoEllipse } from 'react-icons/io5';
import { Trans } from '@lingui/react';

type MembersProps = {
  key?: string;
  name: string;
  isTestNet: boolean;
  isMainNet: boolean;
  memberId: string;
};

const isActivated = (status: boolean) => {
  return status ? '#34A853' : '#EA4335';
};

const CertificateRegistrationRow: FC<MembersProps> = (props) => {
  const { name, isTestNet, isMainNet, memberId } = props;
  const textColor = useColorModeValue('#858585', 'white');
  return (
    <Tr
      border="1px solid #23A7E0"
      borderRadius={100}
      sx={{
        td: {
          height: '66px',
          borderTop: '1px solid #23A7E0',
          borderBottom: '1px solid #23A7E0'
        },
        'td:first-child': {
          border: '1px solid #23A7E0',
          borderLeftRadius: 100,
          borderRight: 'none'
        },
        'td:last-child': {
          border: '1px solid #23A7E0',
          borderRightRadius: 100,
          borderLeft: 'none'
        }
      }}>
      <Td>
        <Flex py=".8rem" minWidth="100%" flexWrap="nowrap" textAlign={'left'} verticalAlign={''}>
          <Text fontSize="md" minWidth="100%">
            {name}
          </Text>
        </Flex>
      </Td>
      <Td>
        <Icon as={IoEllipse} h={30} w={31} color={isActivated(isTestNet)} />
      </Td>
      <Td>
        <Icon as={IoEllipse} h={30} w={31} color={isActivated(isMainNet)} />
      </Td>
      <Td pr={0}>
        <Button
          bg={'#55ACD8'}
          color={'white'}
          _hover={{ bg: '#55ACF8' }}
          _active={{
            bg: '#55ACF8',
            transform: 'scale(0.98)',
            borderColor: '#55ACE8'
          }}
          _focus={{
            boxShadow: '0 0 1px 2px rgba(88, 144, 255, .75), 0 1px 1px rgba(0, 0, 0, .15)'
          }}>
          <Trans id="Details">Details</Trans>
        </Button>
      </Td>
    </Tr>
  );
};

export default CertificateRegistrationRow;
