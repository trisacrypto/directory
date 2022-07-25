import { Stack, Spinner, Flex, Box } from '@chakra-ui/react';

const Loader: React.FC = () => {
  return (
    <Flex
      height={'100vh'}
      bg={'white'}
      alignItems={'center'}
      textAlign={'center'}
      justifyContent={'center'}>
      <Spinner color="blue.500" size="xl" />
    </Flex>
  );
};

export default Loader;
