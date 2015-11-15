FROM node:0.12.7

# Install Badge App
RUN mkdir -p /node-root
RUN git clone https://github.com/badges/shields.git /node-root
WORKDIR /node-root
RUN npm install

# Install Go App
RUN mkdir /go-root
RUN chmod -R 0777 /go-root
WORKDIR /go-root

COPY dist/linux_amd64_gitlab-coverage-badge ./coverage
RUN chmod +x ./coverage
CMD ["./coverage"]

EXPOSE 8080