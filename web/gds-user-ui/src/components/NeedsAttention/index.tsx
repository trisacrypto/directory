import { Box, Text, Stack, Button, HStack } from '@chakra-ui/react';
import FormButton from 'components/ui/FormButton';
import { NavLink } from 'react-router-dom';
const NeedsAttention = () => {
  return (
    <Stack
      minHeight={67}
      bg={'#D8EAF6'}
      p={5}
      border="1px solid #DFE0EB"
      fontSize={18}
      display={'flex'}>
      <HStack justifyContent={'space-between'}>
        <Text fontWeight={'bold'}> Needs Attention </Text>
        <Text> Complete Testnet Registration </Text>

        <Box>
          <NavLink to="/dashboard/certificate/registration">
            <Button
              width={142}
              as={'a'}
              borderRadius={0}
              background="#55ACD8"
              color="#fff"
              _hover={{ background: 'blue.200' }}>
              Start
            </Button>
          </NavLink>
        </Box>
      </HStack>
    </Stack>
  );
};

export default NeedsAttention;
