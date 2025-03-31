// import userEvent from "@testing-library/user-event"
// import faker from "faker"
// import TrisaDetails from "pages/app/details/BasicDetails/components/TrisaDetails"
// import { render, screen, waitFor } from "utils/test-utils"

describe('<TrisaDetails />', () => {
  it.todo('should render data');

  // it("should render data", () => {
  //     const data = {
  //         name: "Guidehouse Inc.",
  //         vasp: {
  //             business_category: "BUSINESS_ENTITY",
  //             common_name: "trisa-axxxxxx16cc67xxxxxxxxxxdb.traveler.ciphertrace.com",
  //             id: "03faf7d2-451d-4d90-8302-e80f0cc9848a",
  //             registered_directory: "testnet.directory",
  //             trisa_endpoint: "trisa-xxxxxxxxx2abexxxx.traveler.ciphertrace.com:443",
  //         },
  //     }

  //     const handleTrisaJsonExportClick = jest.fn()

  //     render(<TrisaDetails data={data} handleTrisaJsonExportClick={handleTrisaJsonExportClick} />)

  //     expect(screen.getByText(/id:/i).firstElementChild.textContent).toBe(data.vasp.id)
  //     expect(screen.getByText(/common name:/i).firstElementChild.textContent).toBe(data.vasp.common_name)
  //     expect(screen.getByText(/endpoint:/i).firstElementChild.textContent).toBe(data.vasp.trisa_endpoint)
  //     expect(screen.getByText(/registered directory:/i).firstElementChild.textContent).toBe(data.vasp.registered_directory)

  // })

  // it("should download JSON file", async () => {
  //     const data = {
  //         name: "Guidehouse Inc.",
  //         vasp: {
  //             common_name: "trisa-axxxxxx16cc67xxxxxxxxxxdb.traveler.ciphertrace.com",
  //             id: faker.datatype.uuid(),
  //             registered_directory: "testnet.directory",
  //             trisa_endpoint: "trisa-xxxxxxxxx2abexxxx.traveler.ciphertrace.com:443",
  //         },
  //     }

  //     const handleTrisaJsonExportClick = jest.fn()

  //     render(<TrisaDetails data={data} handleTrisaJsonExportClick={handleTrisaJsonExportClick} />)

  //     const downloadEl = screen.getByRole('button', { title: /download as json/i })
  //     await waitFor(() => {
  //         userEvent.click(downloadEl)
  //     })

  //     expect(handleTrisaJsonExportClick).toHaveBeenCalled()
  // })
});
