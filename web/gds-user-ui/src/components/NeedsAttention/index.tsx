import { Box, Text, Stack, Button } from '@chakra-ui/react';
import FormButton from 'components/ui/FormButton';

const NeedsAttention = () => {
  return (
    <Stack minHeight={67} bg={'#D8EAF6'} p={5} border="1px solid #DFE0EB" fontSize={18}>
      <Stack direction={'row'} spacing={3} alignItems="center">
        <Stack direction={['column', 'row']} spacing={3}>
          <Text fontWeight={'bold'}> Needs Attention </Text>
          <Text> Complete Testnet Registration </Text>
        </Stack>
        <Box>
          <Button
            width={142}
            as={'a'}
            href={'/dashboard/certificate/registration'}
            borderRadius={0}
            background="#55ACD8"
            color="#fff"
            _hover={{ background: 'blue.200' }}>
            Start
          </Button>
        </Box>
      </Stack>
    </Stack>
  );
};

export default NeedsAttention;
