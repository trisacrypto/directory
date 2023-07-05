import { FormControl, FormLabel, Select } from "@chakra-ui/react";
import { Trans } from "@lingui/macro";
import { useState } from "react";
import { DirectoryTypeEnum } from "../memberType";
import { getMembersService } from "../service";

const MemberSelectNetwork = () => {
  const [selectNetwork, setSelectNetwork] = useState("mainnet");
  getMembersService(selectNetwork as any);

    return (
    <FormControl>
      <FormLabel>
        <Trans>Select Network</Trans>
      </FormLabel>
      <Select
        data-testid="select-network"
        onChange={(e) => setSelectNetwork((e.target.value))}
      >
        <option value={DirectoryTypeEnum.MAINNET}>MainNet</option>
        <option value={DirectoryTypeEnum.TESTNET}>TestNet</option>
      </Select>
    </FormControl>
    );
};

export default MemberSelectNetwork;

