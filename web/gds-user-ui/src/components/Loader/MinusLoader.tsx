import { Spinner, Flex, Text, VStack } from '@chakra-ui/react';

interface LoaderProps {
  text?: string;
}
// This loader is for containers that are loading
const MinusLoader: React.FC<LoaderProps> = (props) => {
  return (
    <Flex bg={'white'} alignItems={'center'} textAlign={'center'} justifyContent={'center'}>
      <VStack spacing={4} m={'auto'} py={10}>
        <Spinner color="blue.500" size="xl" />
        <Text>{props?.text}</Text>
      </VStack>
    </Flex>
  );
};
MinusLoader.defaultProps = {
  text: 'Loading...'
};

export default MinusLoader;
