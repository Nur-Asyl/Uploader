import React, {useState} from "react";

import axios from "axios";

import {Button, Form, Container, Modal} from "react-bootstrap";

const Uploaders = () => {

    const [refreshPage, setRefreshPage] = useState(false);

    const [uploadExcelState, setUploadExcelState] = useState(false);
    const [uploadXMLState, setUploadXMLState] = useState(false);

    const [errorState, setErrorState] = useState({
        response: null,
        affected: null,
        errorMessage: null
    })

    return (
      <div>
          {/* upload excel button */}
          <Container>
              <Button className="btn-info" onClick={() => setUploadExcelState(true)}>Upload Excel</Button>
          </Container>

          {/* upload xml button */}
          <Container>
              <Button className="btn-info" onClick={() => setUploadXMLState(true)}>Upload XML</Button>
          </Container>

          {/* excel modal */}
          <Modal size="lg" show={uploadExcelState} onHide={() => setUploadExcelState(false)} centered>
              <Modal.Header closeButton>
                  <Modal.Title>Upload Excel</Modal.Title>
              </Modal.Header>

              <Modal.Body>
                  <Form className="m-3">
                      <Form.Group className="mb-3">
                          <Form.Label>JSON REQUEST</Form.Label>
                          <Form.Control as="textarea" id="jsonData" rows={10} placeholder="Enter JSON"/>
                      </Form.Group>

                      <Form.Group className="mb-3">
                          <Form.Label>Upload Excel File</Form.Label>
                          <Form.Control type="file" id="formFile"/>
                      </Form.Group>

                      <Button className="float-start btn-info w-25" onClick={() => {uploadExcel()}}>Upload</Button>
                      <Button className="float-end btn-info w-25" onClick={() => setUploadExcelState(false)}>Cancel</Button>
                  </Form>
              </Modal.Body>
          </Modal>

          {/* xml modal */}
          <Modal size="lg" show={uploadXMLState} onHide={() => setUploadXMLState(false)} centered>
              <Modal.Header closeButton>
                  <Modal.Title>Upload XML</Modal.Title>
              </Modal.Header>

              <Modal.Body>
                  <Form className="m-3">
                      <Form.Group className="mb-3">
                          <Form.Label>JSON REQUEST</Form.Label>
                          <Form.Control as="textarea" id="jsonData" rows={10} placeholder="Enter JSON"/>
                      </Form.Group>

                      <Form.Group className="mb-3">
                          <Form.Label>Upload XML File</Form.Label>
                          <Form.Control type="file" id="formFile"/>
                      </Form.Group>

                      <Button className="float-start btn-info w-25" onClick={() => {uploadXML()}}>Upload</Button>
                      <Button className="float-end btn-info w-25" onClick={() => setUploadXMLState(false)}>Cancel</Button>
                  </Form>
              </Modal.Body>
          </Modal>

          <Container className="border rounded-5 border-info mt-5 border-3">
              {errorState.response ? <div className="mb-3">Status: {errorState.response}</div> : null }
              {errorState.affected ?
                  <div>
                      <p>Total: {errorState.affected.total} </p>
                      <p>Inserted: {errorState.affected.inserted} </p>
                      <p>Failed Rows: </p>
                      <p>{errorState.affected.failedRows}</p>
                  </div> : null }
              {errorState.errorMessage ? <div>Error: {errorState.errorMessage}</div> : null }
          </Container>
      </div>
    );

    async function uploadExcel() {
        setUploadExcelState(false)

        // Receiving uploaded file
        let fileInput = document.getElementById("formFile");
        if(fileInput.files.length === 0) {
            setErrorState({
                errorMessage: "No file found. Upload a file"
            })
            return
        }
        let file = fileInput.files[0];

        // Check if file is an Excel file
        const validMimeTypes = [
            "application/vnd.ms-excel",
            "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
        ];
        if (!validMimeTypes.includes(file.type)) {
            setErrorState({
                response: null,
                affected: null,
                errorMessage: "Invalid file type. Require Excel file"
            });
            return;
        }

        // Receiving json from textarea
        let jsonData = document.getElementById("jsonData").value;
        if(jsonData.length === 0) {
            setErrorState({
                errorMessage: "No json found"
            })
            return
        }
        // Creating FormData object
        let formData = new FormData();
        formData.append("json_data", jsonData);

        // Creating blob form with file and it's type
        let fileBlob = new Blob([file], { type: file.type });
        formData.append("file", fileBlob, file.name);

        let url = "http://localhost:8000/upload/excel"

        await axios.post(url, formData, {
            headers: {
                "Content-Type": `multipart/form-data`
            }
        }).then(response => {
            console.log(response.data)
            setErrorState({
                response: response.status,
                affected: {
                    total: response.data.Total,
                    inserted: response.data.Inserted,
                    failedRows: response.data.FailedRows
                },
                errorMessage: null
            })
        }).catch((error) => {
            if (error.response) {
                setErrorState({
                    response: error.response.status,
                    affected: null,
                    errorMessage: error.response.data
                });
            } else {
                setErrorState({
                    response: null,
                    affected: null,
                    errorMessage: error.message
                });
            }

        });
    }

    async function uploadXML() {
        setUploadXMLState(false)

        // Receiving uploaded file
        let fileInput = document.getElementById("formFile");
        if(fileInput.files.length === 0) {
            setErrorState({
                errorMessage: "No file found. Upload a file"
            })
            return
        }
        let file = fileInput.files[0];

        // Check if file is an Excel file
        const validMimeTypes = [
            "text/xml"
        ];
        if (!validMimeTypes.includes(file.type)) {
            setErrorState({
                response: null,
                affected: null,
                errorMessage: "Invalid file type. Require XML file"
            });
            return;
        }

        // Receiving json from textarea
        let jsonData = document.getElementById("jsonData").value;
        if(jsonData.length === 0) {
            setErrorState({
                errorMessage: "No json found"
            })
            return
        }
        // Creating FormData object
        let formData = new FormData();
        formData.append("json_data", jsonData);

        // Creating blob form with file and it's type
        let fileBlob = new Blob([file], { type: file.type });
        formData.append("file", fileBlob, file.name);

        let url = "http://localhost:8000/upload/xml"

        await axios.post(url, formData, {
            headers: {
                "Content-Type": `multipart/form-data`
            }
        }).then(response => {
            console.log(response.data)
            setErrorState({
                response: response.status,
                affected: {
                    total: response.data.Total,
                    inserted: response.data.Inserted,
                    failedRows: response.data.FailedRows
                },
                errorMessage: null
            })
        }).catch((error) => {
            if (error.response) {
                setErrorState({
                    response: error.response.status,
                    affected: null,
                    errorMessage: error.response.data
                });
            } else {
                setErrorState({
                    response: null,
                    affected: null,
                    errorMessage: error.message
                });
            }

        });
    }
}

export default Uploaders;
