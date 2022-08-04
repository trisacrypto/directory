import { Box, Text, Stack, Button, HStack } from '@chakra-ui/react';

import { NavLink } from 'react-router-dom';

import { Trans } from '@lingui/react';

export type NeedsAttentionProps = {
  text: string;
  buttonText: string;
  onClick?: (ev?: any) => void;
};

const NeedsAttention = ({ text, buttonText, onClick }: NeedsAttentionProps) => {
  return (
    <Stack
      minHeight={67}
      bg={'#D8EAF6'}
      p={5}
      border="1px solid #eee"
      fontSize={18}
      display={'flex'}
      borderRadius={'10px'}>
      <HStack justifyContent={'space-between'}>
        <Text fontWeight={'bold'}>
          <Trans id="Needs Attention">Needs Attention</Trans>
        </Text>
        <Text>
          <Trans id="Complete Testnet Registration">Complete Testnet Registration</Trans>
        </Text>

        <Box>
          <Button
            onClick={onClick}
            width={142}
            as={'a'}
            borderRadius={0}
            background="#55ACD8"
            color="#fff"
            cursor="pointer"
            _hover={{ background: 'blue.200' }}>
            {buttonText}
          </Button>
        </Box>
      </HStack>
    </Stack>
  );
};

export default NeedsAttention;
