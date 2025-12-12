float triangle(float x, float period) { 
	return 2.0 * abs(3.0*   ((x / period) - floor((x / period) + 0.5))) - 1.0;
}
 
float field(in vec3 position) {	
  float strength = 7.0 + 0.03 * log(1.0e-6 + fract(sin(iTime) * 373.11));
  float accumulated = 0.0;
  float previousMagnitude = 0.0;
  float totalWeight = 0.0;	

  for (int i = 0; i < 6; ++i) {
    float magnitude = dot(position, position);
    position = abs(position) / magnitude + vec3(-0.5, -0.8 + 0.1 * sin(-iTime * 0.1 + 2.0), -1.1 + 0.3 * cos(iTime * 0.3));
    float weight = exp(-float(i) / 7.0);
    accumulated += weight * exp(-strength * pow(abs(magnitude - previousMagnitude), 2.3));
    totalWeight += weight;
    previousMagnitude = magnitude;
  }

  return max(0.0, 5.0 * accumulated / totalWeight - 0.7);
}

void mainImage(out vec4 fragColor, in vec2 fragCoord) {
  const float baseSpeed = 0.02;
  const int maxIterations = 16;
  const float formulaParameter = 0.79;
  const float volumeSteps = 7.0;
  const float stepSize = 0.24;
  const float zoomFactor = 0.1;
  const float tilingFactor = 0.85;
  const float baseBrightness = 0.0008;
  const float darkMatter = 0.2;
  const float distanceFading = 0.56;
  const float colorSaturation = 0.9;
  const float transverseMotion = 0.2;
  const float cloudOpacity = 0.48;
  const float zoomSpeed = 0.0002;

  vec2 normalizedCoordinates = 2.0 * fragCoord.xy / vec2(512) - 1.0;
  vec2 scaledCoordinates = normalizedCoordinates * vec2(512) / 512.0;

  float timeElapsed = iTime;               
  float speedAdjustment = -baseSpeed;
  float formulaAdjustment = formulaParameter;

  speedAdjustment = zoomSpeed * cos(iTime * 0.02 + 3.1415926 / 4.0);          

  vec2 uvCoordinates = scaledCoordinates;		       

  float rotationXZ = 0.9;
  float rotationYZ = -0.6;
  float rotationXY = 0.9 + iTime * 0.08;	

  mat2 rotationMatrixXZ = mat2(vec2(cos(rotationXZ), sin(rotationXZ)), vec2(-sin(rotationXZ), cos(rotationXZ)));	
  mat2 rotationMatrixYZ = mat2(vec2(cos(rotationYZ), sin(rotationYZ)), vec2(-sin(rotationYZ), cos(rotationYZ)));		
  mat2 rotationMatrixXY = mat2(vec2(cos(rotationXY), sin(rotationXY)), vec2(-sin(rotationXY), cos(rotationXY)));

  vec2 canvasCenter = vec2(0.5, 0.5);
  vec3 rayDirection = vec3(uvCoordinates * zoomFactor, 1.0); 
  vec3 cameraPosition = vec3(0.0, 0.0, 0.0);                               
  cameraPosition.x -= 2.0 * (canvasCenter.x - 0.5);
  cameraPosition.y -= 2.0 * (canvasCenter.y - 0.5);

  vec3 forwardVector = vec3(0.0, 0.0, 1.0);   
  cameraPosition.x += transverseMotion * cos(0.01 * iTime) + 0.001 * iTime;
  cameraPosition.y += transverseMotion * sin(0.01 * iTime) + 0.001 * iTime;
  cameraPosition.z += 0.003 * iTime;	

  rayDirection.xz *= rotationMatrixXZ;
  forwardVector.xz *= rotationMatrixXZ;	
  rayDirection.yz *= rotationMatrixYZ;
  forwardVector.yz *= rotationMatrixYZ;

  cameraPosition.xy *= -1.0 * rotationMatrixXY;
  cameraPosition.xz *= rotationMatrixXZ;
  cameraPosition.yz *= rotationMatrixYZ;

  float zoomOffset = (timeElapsed - 3311.0) * speedAdjustment;
  cameraPosition += forwardVector * zoomOffset;
  float sampleOffset = mod(zoomOffset, stepSize);
  float normalizedSampleOffset = sampleOffset / stepSize;

  float stepDistance = 0.24;
  float secondaryStepDistance = stepDistance + stepSize / 2.0;
  vec3 accumulatedColor = vec3(0.0);
  float fieldContribution = 0.0;	
  vec3 backgroundColor = vec3(0.0);

  for (float stepIndex = 0.0; stepIndex < volumeSteps; ++stepIndex) {
    vec3 primaryPosition = cameraPosition + (stepDistance + sampleOffset) * rayDirection;
    vec3 secondaryPosition = cameraPosition + (secondaryStepDistance + sampleOffset) * rayDirection;

    primaryPosition = abs(vec3(tilingFactor) - mod(primaryPosition, vec3(tilingFactor * 2.0)));
    secondaryPosition = abs(vec3(tilingFactor) - mod(secondaryPosition, vec3(tilingFactor * 2.0)));

    fieldContribution = field(secondaryPosition);

    float particleAccumulator = 0.0, particleDistance = 0.0;
    for (int i = 0; i < maxIterations; ++i) {
      primaryPosition = abs(primaryPosition) / dot(primaryPosition, primaryPosition) - formulaAdjustment;
      float distanceChange = abs(length(primaryPosition) - particleDistance);
      particleAccumulator += i > 2 ? min(12.0, distanceChange) : distanceChange;
      particleDistance = length(primaryPosition);
    }
    particleAccumulator *= particleAccumulator * particleAccumulator;

    float fadeFactor = pow(distanceFading, max(0.0, float(stepIndex) - normalizedSampleOffset));
    accumulatedColor += vec3(stepDistance, stepDistance * stepDistance, stepDistance * stepDistance * stepDistance * stepDistance) 
                        * particleAccumulator * baseBrightness * fadeFactor;
    backgroundColor += mix(0.4, 1.0, cloudOpacity) * vec3(1.8 * fieldContribution * fieldContribution * fieldContribution, 
                                                          1.4 * fieldContribution * fieldContribution, fieldContribution) * fadeFactor;
    stepDistance += stepSize;
    secondaryStepDistance += stepSize;		
  }
  
  accumulatedColor = mix(vec3(length(accumulatedColor)), accumulatedColor, colorSaturation);

  vec4 foregroundColor = vec4(accumulatedColor * 0.01, 1.0);	
  backgroundColor *= cloudOpacity;	
  backgroundColor.b *= 1.8;
  backgroundColor.r *= 0.05;

  backgroundColor.b = 0.5 * mix(backgroundColor.g, backgroundColor.b, 0.8);
  backgroundColor.g = 0.0;
  backgroundColor.bg = mix(backgroundColor.gb, backgroundColor.bg, 0.5 * (cos(iTime * 0.01) + 1.0));	

  vec2 terminalUV = fragCoord.xy / iResolution.xy;
  vec4 terminalColor = texture(iChannel0, terminalUV);

  float brightnessThreshold = 0.1;
  float terminalBrightness = dot(terminalColor.rgb, vec3(0.2126, 0.7152, 0.0722));

  if (terminalBrightness < brightnessThreshold) {
    fragColor = mix(terminalColor, vec4(foregroundColor.rgb + backgroundColor, 1.0), 0.24);
  } else {
    fragColor = terminalColor;
  }
}

