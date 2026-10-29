export async function handleFileInput(event) {
  const imageInfo = event.target.files[0];

  if (imageInfo) {
    if (imageInfo.size > 8000000) {
      alert("Image file size needs to be 8MB or smaller!");
      return { dataUrl: null, imageInfo: null }; // Return an object with null values in case of an error
    }

    const reader = new FileReader();

    return new Promise((resolve, reject) => {
      reader.onload = (event) => {
        const dataUrl = event.target.result;
        resolve({ dataUrl, imageInfo }); // Resolve the promise with an object containing both values
      };

      reader.onerror = (error) => {
        reject(error); // Reject the promise in case of an error
      };

      reader.readAsDataURL(imageInfo);
    });
  }
  return { dataUrl: null, imageInfo: null }; // Return an object with null values if no userAvatar
}