

const bufferBase64 = Buffer.from("3c1a01a1", 'base64');
const hex = bufferBase64.toString('hex');
console.log(hex)

// '387252', '2024-12-28 08:53:51', '2024-12-28 08:53:51', '{\"deviceInfo\":{\"tenantId\":\"52f14cd4-c6f1-4fbd-8f87-4025e1d49242\",\"tenantName\":\"IoT Solutions\",\"applicationId\":\"979fe51e-7732-4da0-82e4-0b0b48acc2d7\",\"applicationName\":\"IoT Park - EU868\",\"deviceProfileId\":\"a66b17ae-711d-47f7-a0b4-7629fca2e94f\",\"deviceProfileName\":\"IoT Park - EU868\",\"deviceName\":\"AC1F09FFFE158666\",\"devEui\":\"ac1f09fffe158666\",\"tags\":{}},\"data\":\"O2dvv+caDQAEAK0A\"}', '3', '1'

