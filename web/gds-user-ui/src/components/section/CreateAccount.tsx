import {
  Flex,
  Box,
  FormControl,
  FormLabel,
  Input,
  Checkbox,
  Stack,
  Link,
  Button,
  Heading,
  Text,
  useColorModeValue,
  Image,
} from "@chakra-ui/react";

import { GoogleIcon } from "components/Icon";

import { colors } from "utils/theme";

export default function CreateAccount() {
  return (
    <Flex
      minH={"100vh"}
      align={"center"}
      justify={"center"}
      fontFamily={"open sans"}
      fontSize={"xl"}
      bg={useColorModeValue("white", "gray.800")}
    >
      <Stack spacing={8} mx={"auto"} maxW={"lg"} py={12} px={6}>
        <Stack align={"center"}>
          <Heading fontSize={"4xl"}></Heading>
          <Text color={useColorModeValue("gray.600", "white")}>
            <Text as={"span"} fontWeight={"bold"}>
              Create your TRISA account.
            </Text>{" "}
            We recommend that a senior compliance officer initialally creates
            the account for the VASP. Additional accounts can be created later.
          </Text>
        </Stack>
        <Stack align={"center"} justify={"center"} fontFamily={"open sans"}>
          <Button
            bg={"gray.100"}
            w="100%"
            height={"64px"}
            color={"gray.600"}
            _hover={{
              background: useColorModeValue("gray.200", "black"),
              color: useColorModeValue("gray.600", "white"),
            }}
          >
            <GoogleIcon h={24} />
            <Text as={"span"} ml={3}>
              Continue with google
            </Text>
          </Button>
          <Text py={3}>Or</Text>
        </Stack>

        <Box
          rounded={"lg"}
          bg={useColorModeValue("white", "transparent")}
          width={"100%"}
          position={"relative"}
          bottom={5}
        >
          <Stack spacing={4}>
            <FormControl id="email">
              <Input type="email" height={"64px"} placeholder="Email Adress" />
            </FormControl>
            <FormControl id="password">
              <Input type="password" height={"64px"} placeholder="Password" />
            </FormControl>
            <Stack spacing={10}>
              <Button
                bg={colors.system.blue}
                color={"white"}
                height={"64px"}
                _hover={{
                  bg: "#10aaed",
                }}
              >
                Create an Account
              </Button>
              <Text textAlign="center">
                Already have an account?{" "}
                <Link href="/login" color={colors.system.cyan}>
                  {" "}
                  Log in.
                </Link>
              </Text>
            </Stack>
          </Stack>
        </Box>
      </Stack>
    </Flex>
  );
}
