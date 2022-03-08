import { Meta, Story } from "@storybook/react";
import NameIdentifiers from ".";

type NameIdentifiersProps = {};

export default {
  title: "components/NameIdentifiers",
  component: NameIdentifiers,
} as Meta<NameIdentifiersProps>;

const Template: Story<NameIdentifiersProps> = (args) => (
  <NameIdentifiers {...args} />
);

export const Default = Template.bind({});
Default.args = {};
