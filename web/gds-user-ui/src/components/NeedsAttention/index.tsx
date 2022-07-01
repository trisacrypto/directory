import { Box, Text, Stack, Button, HStack } from '@chakra-ui/react';
import FormButton from 'components/ui/FormButton';
import { NavLink } from 'react-router-dom';

export type NeedsAttentionProps = {
  text: string;
  buttonText: string;
  onClick?: (ev?: any) => void
};

const NeedsAttention = ({ text, buttonText, onClick }: NeedsAttentionProps) => {
  return (
    <Stack
      minHeight={67}
      bg={'#D8EAF6'}
      p={5}
      border="1px solid #555151D4"
      fontSize={18}
      display={'flex'} borderRadius={'10px'}>
      <HStack justifyContent={'space-between'}>
        <Text fontWeight={'bold'}> Needs Attention </Text>
        <Text> {text} </Text>
        <Box>
            <Button
                onClick={onClick}
              width={142}
              as={'a'}
              borderRadius={0}
              background="#55ACD8"
              color="#fff"
              _hover={{ background: 'blue.200' }}>
              {buttonText}
            </Button>
        </Box>
      </HStack>
    </Stack>
  );
};

export default NeedsAttention;
